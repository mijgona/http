package server

import (
	"log"
	"strconv"
	"strings"
)


func (s *Server) convertPath(requestPath string) (map[string]string, string ){

	var paths []string
	var handlerPath string
	wasModified:=false
	log.Print("request path:" + requestPath)
	s.mu.RLock()
	for path := range s.handlers {
		paths = append(paths, path)
	}
	s.mu.RUnlock()
	
	pathParams := make(map[string]string)
	if requestPath[len(requestPath)-1] == '/' {
		requestPath = requestPath[:len(requestPath)-1]
	}


	parts := strings.Split(requestPath, "/")
	log.Print(parts)	
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	for _, v := range paths {
		curParts := strings.Split(v, "/")
		if curParts[len(curParts)-1] == "" {
			curParts = parts[:len(curParts)-1]
		}
		if len(curParts) != len(parts) {
			continue
		}
		temp := make(map[string]string)
		flag := true
		for i := 0; i < len(curParts); i++ {
			str, ok := getMaket(curParts[i])
			if ok {
				ind := strings.Index(curParts[i], "{"+str+"}")				 
				temp[str] = parts[i][ind:]
			} else if curParts[i] != parts[i] {
				flag = false
			}
			if flag {
				for key, value := range temp{
					pathParams[key] = value
					handlerPath=v
					wasModified=true
				}
			}
		}
	}
	for _, v := range paths {
		str := trimMaket(v)			
		curParts := strings.Split(str, "/")
		if curParts[len(curParts)-1] == "" {
			curParts = parts[:len(curParts)-1]
		}
		
		if len(curParts) != len(parts) {
			continue
		}
		//удаляем лишние символы
		flag := true
		for i := 0; i < len(curParts); i++ {
		if curParts[i] != parts[i] {
			flag = false
		}
			if flag {
				if !wasModified {
					handlerPath=str
					wasModified=true
				}
			}
		}
		ind := strings.Index(v, "{"+str+"}")
		log.Print(ind)

	
	}
	
	return pathParams, handlerPath
}
func trimMaket(str string) (string){
	if len(str) < 2 {
		return ""
	}
		answer:=""
	curParts := strings.Split(str, "/")	
	for _, v := range curParts {
		if v==""{
			continue
		}
    	str,_=getMaket(v)
		str1 := "{"+str+"}"
		str2:=strings.Replace(v, str1, "", -1)
		if str2!=""{
			answer+="/"+str2
		}
	}
	return answer
}


func getMaket(str string) (string, bool) {
	if len(str) < 2 {
		return "", false
	}
	cnt1 := strings.Count(str, "{")
	cnt2 := strings.Count(str, "}")
	if cnt1 != 1 || cnt2 != 1 {
		return "", false
	}
	ind1 := strings.Index(str, "{")
	ind2 := strings.Index(str, "}")
	return str[ind1+1 : ind2], true
}

func (s *Server) Response(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		"\r\n" + body
}
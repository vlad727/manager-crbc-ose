package getsa

import (
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/getsacollect"
	"webapp/home/loggeduser"
)

// GetSa возвращает информацию о пользователе в формате YAML.
// Получаем request, парсим, получаем имя пользователя и его группы
// далее передаем имя пользователя и группы в функцию которая нам вернем
// его namespaces и service accounts, где пользователь является админом или rolebinding создан на группу
// все это форматируем в yaml и выгружаем пользователю на web страницу

func GetSa(w http.ResponseWriter, r *http.Request) {

	// send request to parse and get logged user string
	LoggedUser := loggeduser.LoggedUserRun(r)
	log.Printf("Get username from func LoggedUserRun %s", LoggedUser)

	// get map already sorted and slice for crbcmain, slice will be skipped
	// Получаем карту, отсортированную для вывода
	M3, _ := getsacollect.GetSaCollect(LoggedUser)
	//log.Println(Sl1)                                 // get it from getsa collect slice we don't need

	// Marshal to yaml for out to web page
	// Преобразуем карту в YAML
	yamlFile, err := yaml.Marshal(M3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert to string for struct if you do not convert it will be in bytes
	// Конвертируем YAML в строку
	str := string(yamlFile)

	// parse html template
	// Парсим HTML-шаблон
	t, err := template.ParseFiles("tmpl/getsa.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		return
	}

	// init struct and var
	// Передаем данные в шаблон
	Msg := struct {
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: LoggedUser,
	}

	// execute
	// Выполняем рендеринг шаблона
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}

}

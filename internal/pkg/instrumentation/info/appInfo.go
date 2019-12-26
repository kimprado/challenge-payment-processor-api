package info

import (
	"github.com/prometheus/common/version"
)

//App contém informações da aplicação para instrumentação
type App struct {
	Info    Info
	Version Version
}

//NewApp -
func NewApp() (a App) {
	a = App{}
	a.Info = newInfo()
	a.Version.Version = version.Version
	a.Version.Revision = version.Revision
	a.Version.Branch = version.Branch
	a.Version.BuildUser = version.BuildUser
	a.Version.BuildDate = version.BuildDate
	a.Version.GoVersion = version.GoVersion
	return
}

//Info possui descrição da aplicação.
type Info struct {
	Nome      string `json:"nome" default:"Payment Processor API"`
	Descricao string `json:"descricao" default:"Versão simplificada de um processador de pagamentos"`
}

func newInfo() (i Info) {
	i = Info{
		Nome:      "Payment Processor API",
		Descricao: "Versão simplificada de um processador de pagamentos",
	}
	return
}

//Version possui informações sobre versão da aplicação.
//Eventualmente Versao conterá tag da versão em produção
type Version struct {
	VersaoModulo string `json:"versao-modulo"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	Branch       string `json:"branch"`
	BuildUser    string `json:"buildUser"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
}

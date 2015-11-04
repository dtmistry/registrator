package bigip

import (
	"log"
	"net/url"
	"net/http"
	"bytes"
	"crypto/tls"
	"errors"
	"encoding/json"
	"strconv"

	"github.com/dtmistry/registrator/bridge"
)

func init() {
	bridge.Register(new(Factory), "bigip")
}

type Member struct {
	Name string `json:"name"`
	Address string `json:"address"`
}

type BigIpAdapter struct {
	user string
	pass string
	path string
	pool string
	client *http.Client
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	usage := "bigip: invalid bigip URI [%s] e.g.: bigip://<user>:<pass>@<host>/<pool-name>"
	userinfo := uri.User
	if userinfo == nil {
		log.Fatalf(usage, uri)
	}
	user := userinfo.Username()
	pass, _ := userinfo.Password()
	if pass == "" || user == "" {
		log.Fatalf(usage, uri)
	}
	var buffer bytes.Buffer	
	if uri.Host != "" && uri.Path != "" {
		buffer.WriteString("http://")
		buffer.WriteString(uri.Host)
		buffer.WriteString("/mgmt/tm/ltm/pool")
	} else {
		log.Fatalf(usage, uri)
	}
	log.Print("Creating BigIp backend using url : ", buffer.String())
	//Disabling certificate authority check
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	return &BigIpAdapter{user: user, pass: pass, 
		path: buffer.String(), pool: uri.Path, client: &http.Client{Transport: tr}}
}


// Ping will try to connect to consul by attempting to retrieve the current leader.
func (r *BigIpAdapter) Ping() error {
	req, err := http.NewRequest("GET", r.path + r.pool + "/stats", nil)
	req.SetBasicAuth(r.user, r.pass)
	resp, err := r.client.Do(req)
	if(err != nil) {
		log.Println("Err : ", err)
		return err
	}
	if resp.StatusCode != 200 {
		log.Println("Ping response :", resp.Status)
		return errors.New("BigIp pool ["+r.pool+"] is unavailable")
	}
	defer resp.Body.Close()
	return nil
}

func (r *BigIpAdapter) Register(service *bridge.Service) error {

	member := &Member{
		Name: service.Name + ":" + strconv.Itoa(service.Port), 
		Address: service.IP,
	}

	payload, err := json.Marshal(member)
	if err != nil {
		return err;
	}
	req, err := http.NewRequest("POST", r.path + r.pool + "/members", bytes.NewBuffer(payload))
	req.SetBasicAuth(r.user, r.pass)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Unable to add member to BipIp pool ["+r.pool+"]. Status ["+resp.Status+"]")
	}
	defer resp.Body.Close()
	return nil
}


func (r *BigIpAdapter) Deregister(service *bridge.Service) error {

	member := service.Name + ":" + strconv.Itoa(service.Port)

	req, err := http.NewRequest("DELETE", r.path + r.pool + "/members/" + member, nil)
	req.SetBasicAuth(r.user, r.pass)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Unable to delete member from BipIp pool ["+r.pool+"]. Status ["+resp.Status+"]")
	}
	defer resp.Body.Close()
	return nil
}

func (r *BigIpAdapter) Refresh(service *bridge.Service) error {
	return nil
}

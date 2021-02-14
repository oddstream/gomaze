// +build browser

package maze

import (
	"encoding/json"
	"log"
	"runtime"
	"syscall/js"
)

// UserData contains the level the user is on
type UserData struct {
	Copyright       string // Capitals to emit to json
	CompletedLevels int    // Capitals to emit to json
}

// NewUserData create a new UserData object and tries to load it's content from file
// it always returns an object, even if file does not exist
func NewUserData() *UserData {

	ud := &UserData{Copyright: "Copyright ©️ 2021 oddstream.games"}
	// let CompletedLevels default to zero

	if runtime.GOARCH != "wasm" {
		log.Fatal("WASM required")
	}
	println("NewUserData: runtime.GOARCH == WASM")
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	v := localStorage.Get("gomaze")
	// println(v.String())
	bytes := []byte(v.String())
	if len(bytes) > 0 {
		err := json.Unmarshal(bytes[:], ud)
		if err != nil {
			log.Fatal(err)
		}
	}
	return ud
}

// Save writes the UserData object to file
func (ud *UserData) Save() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("WASM required")
	}

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Fatal(err)
	}
	str := string(bytes[:])
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	// localStorage.Set("CompletedLevels", strconv.Itoa(ud.CompletedLevels))
	localStorage.Set("gomaze", str)
}

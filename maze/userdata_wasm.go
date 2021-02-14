// +build browser

package maze

import (
	"log"
	"runtime"
	"strconv"
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
	v := localStorage.Get("CompletedLevels")
	println(v.String())
	i, err := strconv.Atoi(v.String())
	if err == nil {
		ud.CompletedLevels = i
	} else {
		print("error", v.String())
	}

	// globalObject := js.Global()
	// loc := globalObject.Get("location")
	// hrf := loc.Get("href")
	// println(hrf.String())
	return ud
}

// Save writes the UserData object to file
func (ud *UserData) Save() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("WASM required")
	}
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	localStorage.Set("CompletedLevels", strconv.Itoa(ud.CompletedLevels))
}

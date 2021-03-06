package uinput

import (
	"fmt"
	"testing"
)

// This test will confirm that basic key events are working.
// Note that only Key1 is used here, as the purpose of this test is to ensure that the event handling for
// keyboard devices is working. All other keys, defined in keycodes.go should work as well if this test passes.
// Another thing to keep in mind is that there are certain key codes that might not be great candidates for
// unit testing, as they may create unwanted side effects, like logging out the current user, etc...
func TestKeysInValidRangeWork(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}

	err = vk.KeyPress(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(keyReserved)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
	}

	err = vk.KeyPress(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key press. Last error was: %s\n", err)
	}

	err = vk.KeyDown(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key down event. Last error was: %s\n", err)
	}

	err = vk.KeyUp(keyMax)
	if err != nil {
		t.Fatalf("Failed to send key up event. Last error was: %s\n", err)
	}

	err = vk.Close()

	if err != nil {
		t.Fatalf("Failed to close device. Last error was: %s\n", err)
	}
}

func TestKeyboardCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	CreateKeyboard("", []byte("KeyboardDevice"))
	t.Fatalf("Empty path did not yield a panic")
}

func TestKeyboardCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"
	expected := "device path '" + path + "' does not exist"
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	CreateKeyboard(path, []byte("KeyboardDevice"))
	t.Fatalf("Invalid path did not yield a panic")
}

func TestKeyboardCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	defer func() {
		if r := recover(); r != nil {
			actual := r.(string)
			if actual != expected {
				t.Fatalf("Expected: %s\nActual: %s", expected, actual)
			}
		}
	}()
	CreateKeyboard("/dev/uinput", []byte(name))
	t.Fatalf("Invalid name did not yield a panic")
}

func TestKeyOutsideOfRangeKeyPressFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyPress(249)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyPress(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}
func TestKeyOutsideOfRangeKeyUpFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyUp(249)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyUp(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}

func TestKeyOutsideOfRangeKeyDownFails(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	defer vk.Close()

	err = vk.KeyDown(249)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

	err = vk.KeyDown(-1)
	if err == nil {
		t.Fatalf("Expected key press to fail due to invalid key code, but got no error.")
	}

}

func TestKeyPressFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyPress(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}

func TestKeyUpFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyUp(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}

func TestKeyDownFailsIfDeviceIsClosed(t *testing.T) {
	vk, err := CreateKeyboard("/dev/uinput", []byte("Test Basic Keyboard"))
	if err != nil {
		t.Fatalf("Failed to create the virtual keyboard. Last error was: %s\n", err)
	}
	vk.Close()

	err = vk.KeyDown(Key1)
	if err == nil {
		t.Fatalf("Expected KeyPress to fail, but no error was returned.")
	}
}

Uinput [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Build Status](https://travis-ci.org/bendahl/uinput.svg?branch=master)](https://travis-ci.org/bendahl/uinput) [![GoDoc](https://godoc.org/github.com/bendahl/uinput?status.png)](https://godoc.org/github.com/bendahl/uinput) [![Go Report Card](https://goreportcard.com/badge/github.com/bendahl/uinput)](https://goreportcard.com/report/github.com/bendahl/uinput)
====

This package provides pure go wrapper functions for the LINUX uinput device, which allows to create virtual input devices 
in userspace. At the moment this package offers a virtual keyboard implementation as well as a virtual mouse device and
a touch pad device. 

The keyboard can be used to either send single key presses or hold down a specified key and release it later 
(useful for building game controllers). The mouse device issues relative positional change events to the x and y axis 
of the mouse pointer and may also fire click events (left and right click). For implementing things like region selects
via a virtual mouse pointer, press and release functions for the mouse device are also included.

The touch pad, on the other hand can be used to move the mouse cursor to the specified position on the screen and to
issue left and right clicks. Note that you'll need to specify the region size of your screen first though (happens during
device creation).

Please note that you will need to make sure to have the necessary rights to write to uinput. You can either chmod your 
uinput device, or add a rule in /etc/udev/rules.d to allow your user's group or a dedicated group to write to the device.
You may use the following two commands to add the necessary rights for you current user to a file called 99-$USER.rules 
(where $USER is your current user's name):
<pre><code>
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
</code></pre>

Installation
-------------
Simply check out the repository and use the commands <pre><code>go build && go install</code></pre> 
The package will then be installed to your local respository, along with the package documentation. 
The documentation contains more details on the usage of this package. 


Alternatively, if you'd like to use a specific version/tag of this library, you may use the gopkg.in service.
To get v.1.x.x (latest stable v1), simply enter this on your command line:

<code>go get gopkg.in/bendahl/uinput.v1</code>

To import this version in your code use:

<code>import "gopkg.in/bendahl/uinput.v1"</code>

You may then refer to it as "uinput" in your code. For further details see: https://gopkg.in/bendahl/uinput.v1

Thanks to gopkg.in for offering this great service!

Usage
-----
The following section explains some common ways to use this lib.


### Using the virtual keyboard device: 

```go
package main 

import "github.com/bendahl/uinput"
// alternatively (to use specific version), use this:
//import "gopkg.in/bendahl/uinput.v1"

func main() {
	// initialize keyboard and check for possible errors
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer keyboard.Close()

	// prints "a"
	keyboard.KeyPress(uinput.KeyA)

	// prints "A"
	// Note that you could use caps lock instead of using shift with KeyDown and KeyUp
	keyboard.KeyDown(uinput.KeyLeftshift)
	keyboard.KeyPress(uinput.KeyA)
	keyboard.KeyUp(uinput.KeyLeftshift)

	// prints "00000"
	for i := 0; i < 5; i++ {
		keyboard.KeyPress(uinput.Key0)
	}
}
```

### Using the virtual mouse device:

```go
package main

import "github.com/bendahl/uinput"
// alternatively (to use specific version), use this:
//import "gopkg.in/bendahl/uinput.v1"

func main() {
	// initialize mouse and check for possible errors
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer mouse.Close()

	// mouse pointer will be moved up by 10 pixels
	mouse.MoveUp(10)
	// mouse pointer will be moved to the right by 10 pixels
	mouse.MoveRight(10)
	// mouse pointer will be moved down by 10 pixels
	mouse.MoveDown(10)
	// mouse pointer will be moved to the left by 10 pixels (we're back to where we started)
	mouse.MoveLeft(10)

	// click left
	mouse.LeftClick()
	// click right (depending on context a context menu may appear)
	mouse.RightClick()

	// hold down left mouse button
	mouse.LeftPress()
	// move mouse pointer down by 100 pixels while holding down the left key
	mouse.MoveDown(100)
	// release the left mouse button
	mouse.LeftRelease()
}
```

### Using the virtual touch pad device:

```go
package main

import "github.com/bendahl/uinput"
// alternatively (to use specific version), use this:
//import "gopkg.in/bendahl/uinput.v1"

func main() {
	// initialization of the touch device requires to set the screen boundaries
	// min and max values for x and y axis need to be set (usually, 0 should be a sane lower bound)
	touch, err := uinput.CreateTouchPad("/dev/uinput", []byte("testpad"), 0, 800, 0, 600)
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer touch.Close()

	// move pointer to the position 300, 200
	touch.MoveTo(300, 200)
	// press the left mouse key, holding it down
	touch.LeftPress()
	// move pointer to position 400, 400
	touch.MoveTo(400, 400)
	// release the left mouse key
	touch.LeftRelease()
}
```

License
--------
The package falls under the MIT license. Please see the "LICENSE" file for details.

Design Decisions
----------------
After some consideration the validation logic of the uinput module was adapted. Instead of simply returning and error code
it actually panics. The reasoning behind this is that invalid initialization parameters for the uinput device may be a clue
for a programming error that should never pass silently. Also, faulty initialization values would simply cause the uinput
device not to work. "Effective Go" also suggests to use panics in occasions where a library simply cannot set itself up due
 to improper initialization (see: [Effective Go](https://golang.org/doc/effective_go.html#panic)).

Current Status
--------------
2018-03-31: I am happy to announce that v1.0.0 is finally out! Go ahead and use this library in your own projects! Feedback is always welcome.

2019-03-24: Release v1.0.1 fixes a positioning issue that affects the touchpad. See issue #11 for details (positioning works now, but a (possibly) better solution is under investigation).

TODO
----
The current API can be considered stable and the overall functionality (as originally envisioned) is complete. 
Testing on x86_64 and ARM platforms (specifically the RaspberryPi) has been successful. If you'd like to use this library
on a different platform that supports Linux, feel free to test it and share the results. This would be greatly appreciated.
One thing that I'd still like to improve, however, are the test cases. The basic functionality is covered, but more extensive
testing is something that needs to be worked on. 

- [x] Create Tests for the uinput package
- [x] Migrate code from C to GO
- [x] Implement relative input
- [x] Implement absolute input
- [x] Test on different platforms
    - [x] x86_64
    - [x] ARMv6 (RaspberryPi)
- [x] Implement functions to allow mouse button up and down events (for region selects)
- [ ] Extend test cases


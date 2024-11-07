package truncate

import (
	"fmt"
)

func ExampleSmartTruncate() {
	truncateLimit := 20
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaa", []Match{}, truncateLimit, 5))                                                                         // no truncate
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaab", []Match{}, truncateLimit, 5))                                                                        // truncates last char
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaab", []Match{{Offset: 20, Data: []byte("b")}}, truncateLimit, 5))                                         // now we want to see the b
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaac", []Match{{Offset: 20, Data: []byte("b")}}, truncateLimit, 5))                                    // now we want to see the b
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaaac", []Match{{Offset: 20, Data: []byte("b")}}, truncateLimit, 5))                                   // now we want to see the b
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaaaac", []Match{{Offset: 20, Data: []byte("b")}}, truncateLimit, 5))                                  // now we want to see the b
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaac", []Match{{Offset: 20, Data: []byte("b")}, {Offset: 25, Data: []byte("c")}}, truncateLimit, 5))   // now we want to see the b and c
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaaac", []Match{{Offset: 20, Data: []byte("b")}, {Offset: 26, Data: []byte("c")}}, truncateLimit, 5))  // now we want to see the b and c
	fmt.Println(SmartTruncate("aaaaaaaaaaaaaaaaaaaabaaaaaac", []Match{{Offset: 20, Data: []byte("b")}, {Offset: 27, Data: []byte("c")}}, truncateLimit, 5)) // now we want to see the b and c
	fmt.Println(SmartTruncate("aaabaaaaaaaaaaacaaaaa", []Match{{Offset: 3, Data: []byte("b")}, {Offset: 15, Data: []byte("c")}}, truncateLimit, 5))         // now we want to see the b and c
	fmt.Println(SmartTruncate("aaabaaaaaaaaaaaaaaaaa", []Match{{Offset: 3, Data: []byte("b")}}, truncateLimit, 5))                                          // now we want to see the b
	// Output:
	// aaaaaaaaaaaaaaaaaaaa
	// aaaaaaaa[...]aaaaaab
	// aaaaaaaaa[...]aaaaab
	// aaaa[...]aaaaabaaaac
	// aaa[...]aaaaabaaaaac
	// aa[...]aaaaabaaaaaac
	// aaaa[...]aaaaabaaaac
	// aaa[...]aaaaabaaaaac
	// aa[...]aaaaabaaaaaac
	// aaabaaaa[...]acaaaaa
	// aaabaaaaaaaa[...]aaa
}

func ExampleTruncateWithNewlines() {
	fmt.Printf("%q\n", TruncateWithNewlines("aaaaaaaaaaa\naaaaaaaaa", []Match{}, 0, 10))                                                                           // no truncate
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 0, Data: []byte("aaab")}}, 0, 10))                                          // should print first aaab
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 5, Data: []byte("bbbc")}, {Offset: 10, Data: []byte("cccd")}}, 0, 10))      // test with multiple wanted strings
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 0, Data: []byte("aaab\nbbbc")}}, 0, 10))                                    // test multiline wanted string
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 0, Data: []byte("aaab\nbbbc")}, {Offset: 5, Data: []byte("bbbc")}}, 0, 10)) // test multiple strings where one contains the other
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 0, Data: []byte("aaab\nbbbc\ncccd\naaab")}}, 0, 10))                        // test fully wanted string
	fmt.Printf("%q\n", TruncateWithNewlines("aaab", []Match{{Offset: 0, Data: []byte("aaab")}}, 0, 10))                                                            // test without any newlines                                                                // no truncate
	fmt.Printf("%q\n", TruncateWithNewlines("aaab\nbbbc\ncccd\naaab", []Match{{Offset: 0, Data: []byte("aa")}, {Offset: 15, Data: []byte("aa")}}, 0, 10))          // should print both aaab's
	fmt.Printf("%q\n", TruncateWithNewlines("aaaaaaaaaaaaaaabbbbbbbbbbbbbbbbb\nbbbc\ncccd\naaab", []Match{{Offset: 12, Data: []byte("aaab")}}, 15, 5))
	fmt.Printf("%q\n", TruncateWithNewlines("aaaaaaaaaaaaaaabbbbbbbbbbbbbbbbb\nbbbc\ncccd\naaab", []Match{{Offset: 12, Data: []byte("aaab")}}, 10, 5))
	// Output:
	// "aaaaaaaaaaa[...]aaaaaaaaa"
	// "aaab"
	// "bbbc[...]cccd"
	// "aaab[...]bbbc"
	// "aaab[...]bbbc"
	// "aaab[...]bbbc[...]cccd[...]aaab"
	// "aaab"
	// "aaab[...]aaab"
	// "[...]aaabb[...]"
	// "aaa[...]bb"
}

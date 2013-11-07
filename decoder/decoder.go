package decoder

import (
    "fmt"
	"unicode/utf8"
	"unicode/utf16"
	"errors"
)

const (
	E_UNKNOWN byte = iota
	E_UTF8
	E_UTF16_BE
	E_UTF16_LE
	E_UTF32_BE
	E_UTF32_LE
)

var signatures = map[byte] []byte {
	E_UNKNOWN: []byte{},
	E_UTF8: []byte{0xEF, 0xBB, 0xBF},
	E_UTF16_BE: []byte{0xFE, 0xFF},
	E_UTF16_LE: []byte{0xFF, 0xFE},
	E_UTF32_BE: []byte{0x00, 0x00, 0xFE, 0xFF},
	E_UTF32_LE: []byte{0xFF, 0xFE, 0x00, 0x00}}

func print(x interface{}) {
    fmt.Printf("%v", x)
}

func Decode (_input chan byte, output chan rune, useEncoding byte) {
    
    input := make(chan byte)
    
    if E_UNKNOWN == useEncoding {
        //пытаемся узнать кодировку
        go encoding(_input, input)        
	
	    if enc, ok := <-input; ok {
            useEncoding = enc
        } else {
	        //error
	        close(output)
	        return
        }
    } else {
        input = _input
    }
	
	switch useEncoding {
		case E_UNKNOWN:
			//В случае E_UNKNOWN текст декодируется как UTF8
			go decodeUTF8(input, output)
			break
		case E_UTF8:
			go decodeUTF8(input, output)
			break
		case E_UTF16_BE:
			go decodeUTF16(input, output, true)
			break
		case E_UTF16_LE:	
			go decodeUTF16(input, output, false)
			break
		case E_UTF32_BE:
			go decodeUTF32(input, output, true)
			break
		case E_UTF32_LE:
			go decodeUTF32(input, output, false)
			break
		default:
            //error
            //swallow all input than close
			for _ = range input {
			}
			
			close(output)
			break
	}	
}

func encoding(input chan byte, output chan byte) {
	defer close(output)

    //создаем буффер и пытаемся читать в него до 4 байт
    buffer := []byte{}	
    count := 0
    
	for b := range input {
	    
	    if count < 4 {
    	    buffer = append(buffer, b)
	    }
	    
	    count ++
	    	    
	    if 4 == count {	        
	        //проверка сигнатур каждой кодировки
	        flags := make(map[byte] bool, len(signatures))
	        for key, signature := range signatures {
		        flags[key] = len(signature) <= len(buffer)
		
		        for pos, sign := range signature {
			        flags[key] = flags[key] && (sign == buffer[pos])
		        }
	        }
	
	        //обнаруживаем распознанную кодировку. В крайнем случае распознается E_UNKNOWN. 
	        var pos int
	        for i := byte(len(flags) - 1); i >= 0; i-- {
		        if flags[i] {
			        pos = len(signatures[i])
			        //Первый байт = номер кодировки
			        output <- i
			        break
		        }
	        }
	        
	        //кидаем в выходной поток сначала неиспользованные байты буфера
	        for i := pos; i < 4; i++ {
		        output <- buffer[i]
	        }
        }	    

        if count > 4 {
            output <- b
        }
    }
}

func decodeUTF8 (input chan byte, output chan rune) {
	defer close(output)
		
	encodedRune := []byte{}
	for b := range input {
	    
	    if utf8.RuneStart(b) == (len(encodedRune) != 0) {
	        if len(encodedRune) == 0 {
                panic("UTF8 ERROR: Ожидался первый байт цепочки кодированного символа\n")
            } else {
                panic("UTF8 ERROR: Цепочка байт, соответствующая одной кодовой точке, не была декодирована\n")
            }
	    }
	
	    encodedRune = append(encodedRune, b)
	    
	    if utf8.FullRune(encodedRune) {
	        r, _ := utf8.DecodeRune(encodedRune)
	        output <- r
	        
	        encodedRune = []byte{}
	    }
	}
}

func decodeUTF16  (input chan byte, output chan rune, BE bool) {
	defer close(output)
	
	var surrogatePair rune	
	for  {
        r, ok, e := readNBytes(input, 4, BE)
        
        if !ok {
            if nil != e {
                panic(e)
            }            
            break
        }
	    
	    if 0 != surrogatePair {
			if !utf16.IsSurrogate(r) {
				break
				//ошибка поймается в конце
			}		
			r = utf16.DecodeRune(surrogatePair, r)
			
			surrogatePair = 0
		} else {			
			if utf16.IsSurrogate(r) {
				surrogatePair = r
				continue
			}
		}
		
		output <- r
	}
			
	if 0 != surrogatePair {
		panic("UTF16 ERROR: Не хватает второй половины суррогатной пары\n")
	}
}

func decodeUTF32 (input chan byte, output chan rune, BE bool) {
	defer close(output)
	
    for {
        r, ok, e := readNBytes(input, 4, BE)
        
        if !ok {
            if nil != e {
                panic(e)
            }            
            break
        }
        
        output <- r
    }
}

func readNBytes (input chan byte, size int, BE bool) (r rune, ok bool, e error) {
    var (
        t uint32 = 0
        b byte
    )
    
	for i := 0; i < size; i++ {
		b, ok = <-input
				
		//входящий поток завершился
		if !ok {
			if 0 != i {
				e = errors.New("readNBytes ERROR: Количество байт во входном потоке не кратно указанному размеру\n")
			}
			break
		}
				
		if BE {
			t = uint32(t) << uint32(8)
			t += uint32(b)
		} else {			
			t += uint32(b) << uint32(8 * i)
		}
	}
	
	return rune(t), ok, e
}

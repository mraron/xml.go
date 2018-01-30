package xmlgo

import (
   "bufio"
   "errors"
   "io"
   "unicode"
   "fmt"
)

type Attribute struct {
   Key string
   Val string
}

type XMLTag struct {
   Name string
   Attributes []Attribute
   Single bool
   Closing bool
}

type XMLDocument struct {
   Tags []XMLTag
}

func readString(rd *bufio.Reader) ([]rune, error) {
   res := make([]rune, 0)
   for {
      r, _, err := rd.ReadRune()
      fmt.Println(string(r), "readString-ben")

      if err != nil {
         return nil, err
      }

      if !unicode.IsLetter(r) {
         rd.UnreadRune()
         return res, nil
      }else {
         res = append(res, r)
      }
   }

   return nil, nil
}

func readAttr(rd *bufio.Reader) (attr Attribute, err error) {
   r := rune('a')
   state := 0
   for {
      r, _, err = rd.ReadRune()
      if err != nil {
         return
      }
      fmt.Println(string(r), state, "ReadAttr-ban")


      if state == 0 {
         if unicode.IsLetter(r) {
            rd.UnreadRune()

            str := make([]rune, 0)

            str, err = readString(rd)
            attr.Key = string(str)

            state ++
         }else {
            continue
         }
      }else if state == 1 {
         if unicode.IsSpace(r) {
            continue
         }else if r=='=' {
            state ++
         } else {
            err = errors.New("Syntax error 3")
            return
         }
      }else if state == 2 {
         if unicode.IsSpace(r) {
            continue
         } else if r=='"' {
            k := make([]byte, 0)
            k, err = rd.ReadBytes('"')
            if err != nil {
               return
            }

            attr.Val = string(k)
            attr.Val = attr.Val[:len(attr.Val)-1]
            return
         }
      }
   }
}

func Parse(r io.Reader) (*XMLDocument, error) {
   rd := bufio.NewReader(r)
   state := 0

   doc := &XMLDocument{}
   doc.Tags = make([]XMLTag, 0)

   aktTag := XMLTag{}

   for {
      r, _, err := rd.ReadRune()
      if err != nil {
         return doc, nil
      }
      fmt.Println(state, string(r))
      if state == 0 {
         if unicode.IsSpace(r) {
            continue
         }else if r != '<' {
            return nil, errors.New("syntax error 1")
         }else {
            state ++
            aktTag = XMLTag{}
         }
      }else if state == 1 {
         name := ""
         if r == '/' {
            aktTag.Closing = true
         }else {
            rd.UnreadRune()
         }


         rest, err := readString(rd)
         if err != nil {
            return nil, err
         }

         name = string(rest)

         aktTag.Name = name

         state ++
      }else if state==2 {
         if unicode.IsSpace(r) {
            continue
         }else if r=='>' {
            doc.Tags = append(doc.Tags, aktTag)
            state=0
         }else {
            if r=='/' {
               aktTag.Single = true
            }else if unicode.IsLetter(r) {
               rd.UnreadRune();
               attr, err := readAttr(rd)
               if err != nil {
                  return nil, err
               }

               aktTag.Attributes = append(aktTag.Attributes, attr)
            }else {
               return nil, errors.New("syntac error 2")
            }
         }
      }
   }

}

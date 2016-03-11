package filter

import (
   "crypto/md5"
   "crypto/sha1"
   "encoding/hex"
   "fmt"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
   "github.com/rapid7/godap/util"
   "regexp"
   "strings"
)

/////////////////////////////////////////////////
// select filter
/////////////////////////////////////////////////
type FilterSelect struct {
   BaseFilter
}

func (fs *FilterSelect) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   ndoc := make(map[string]interface{})
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         ndoc[k] = docv
      }
   }
   ndocs := make([]map[string]interface{}, 0)
   if len(ndoc) > 0 {
      ndocs = append(ndocs, ndoc)
   }
   return ndocs, nil
}

func init() {
   factory.RegisterFilter("select", func(args []string) (lines api.Filter, err error) {
      filterSelect := &FilterSelect{}
      filterSelect.ParseOpts(args)
      return filterSelect, nil
   })
}

/////////////////////////////////////////////////
// rename filter
/////////////////////////////////////////////////

type FilterRename struct {
   BaseFilter
}

func (fs *FilterRename) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if _, ok := doc[k]; ok {
         doc[v] = doc[k]
         delete(doc, k)
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("rename", func(args []string) (lines api.Filter, err error) {
      filterRename := &FilterRename{}
      filterRename.ParseOpts(args)
      return filterRename, nil
   })
}

/////////////////////////////////////////////////
// remove filter
/////////////////////////////////////////////////

type FilterRemove struct {
   BaseFilter
}

func (fs *FilterRemove) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if _, ok := doc[k]; ok {
         delete(doc, k)
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("remove", func(args []string) (lines api.Filter, err error) {
      filterRemove := &FilterRemove{}
      filterRemove.ParseOpts(args)
      return filterRemove, nil
   })
}

/////////////////////////////////////////////////
// insert filter
/////////////////////////////////////////////////

type FilterInsert struct {
   BaseFilter
}

func (fs *FilterInsert) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      doc[k] = v
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("insert", func(args []string) (lines api.Filter, err error) {
      filterInsert := &FilterInsert{}
      filterInsert.ParseOpts(args)
      return filterInsert, nil
   })
}

/////////////////////////////////////////////////
// include filter
/////////////////////////////////////////////////

type FilterInclude struct {
   BaseFilter
}

func (fs *FilterInclude) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if docv, ok := doc[k]; ok && strings.Contains(docv.(string), v) {
         return []map[string]interface{}{doc}, nil
      }
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("include", func(args []string) (lines api.Filter, err error) {
      filterInclude := &FilterInclude{}
      filterInclude.ParseOpts(args)
      return filterInclude, nil
   })
}

/////////////////////////////////////////////////
// exclude filter
/////////////////////////////////////////////////

type FilterExclude struct {
   BaseFilter
}

func (fs *FilterExclude) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if docv, ok := doc[k]; ok && strings.Contains(docv.(string), v) {
         return make([]map[string]interface{}, 0), nil
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("exclude", func(args []string) (lines api.Filter, err error) {
      filterExclude := &FilterExclude{}
      filterExclude.ParseOpts(args)
      return filterExclude, nil
   })
}

/////////////////////////////////////////////////
// exists filter
/////////////////////////////////////////////////

type FilterExists struct {
   BaseFilter
}

func (fs *FilterExists) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok && len(docv.(string)) > 0 {
         return []map[string]interface{}{doc}, nil
      }
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("exists", func(args []string) (lines api.Filter, err error) {
      filterExists := &FilterExists{}
      filterExists.ParseOpts(args)
      return filterExists, nil
   })
}

/////////////////////////////////////////////////
// notexists filter
/////////////////////////////////////////////////

type FilterNotExists struct {
   BaseFilter
}

func (fs *FilterNotExists) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok && len(docv.(string)) > 0 {
         return make([]map[string]interface{}, 0), nil
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("not_exists", func(args []string) (lines api.Filter, err error) {
      filterNotExists := &FilterNotExists{}
      filterNotExists.ParseOpts(args)
      return filterNotExists, nil
   })
}

/////////////////////////////////////////////////
// where filter
/////////////////////////////////////////////////

type FilterWhere struct {
   key      string
   operator func(string, string) bool
   value    string
   BaseFilter
}

func (fs *FilterWhere) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   if docv, ok := doc[fs.key]; ok && fs.operator(fs.value, docv.(string)) {
      return []map[string]interface{}{doc}, nil
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("where", func(args []string) (lines api.Filter, err error) {
      filterWhere := &FilterWhere{}
      if len(args) != 3 {
         panic(fmt.Sprintf("Expected 3 arguments to 'where' but got %d: %s", len(args), args))
      }
      filterWhere.key = args[0]
      filterWhere.value = args[2]
      if args[1] == "==" {
         filterWhere.operator = func(lhs string, rhs string) bool {
            return lhs == rhs
         }
      } else if args[1] == "!=" {
         filterWhere.operator = func(lhs string, rhs string) bool {
            return lhs != rhs
         }
      } else {
         panic(fmt.Sprintf("Unknown conditional operator for 'where': %s", args[1]))
      }
      return filterWhere, nil
   })
}

/////////////////////////////////////////////////
// split word filter
/////////////////////////////////////////////////

type FilterSplitWord struct {
   regex *regexp.Regexp
   BaseFilter
}

func (fs *FilterSplitWord) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   var lines []map[string]interface{}
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := fs.regex.Split(docv.(string), -1)
         for idx := 0; idx < len(words); idx++ {
            newmap := util.Merge(make(map[string]interface{}), doc)
            newmap[fmt.Sprintf("%s.word", k)] = words[idx]
            lines = append(lines, newmap)
         }
      }
   }
   if len(lines) < 1 {
      return []map[string]interface{}{doc}, nil
   }
   return lines, nil
}

func init() {
   factory.RegisterFilter("split_word", func(args []string) (lines api.Filter, err error) {
      filterSplitWord := &FilterSplitWord{}
      filterSplitWord.regex = regexp.MustCompile("\\W")
      filterSplitWord.ParseOpts(args)
      return filterSplitWord, nil
   })
}

/////////////////////////////////////////////////
// split tab filter
/////////////////////////////////////////////////

type FilterSplitTab struct {
   BaseFilter
}

func (fs *FilterSplitTab) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   var lines []map[string]interface{}
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := strings.Split(docv.(string), ",\t")
         for idx := 0; idx < len(words); idx++ {
            newmap := util.Merge(make(map[string]interface{}), doc)
            newmap[fmt.Sprintf("%s.tab", k)] = words[idx]
            lines = append(lines, newmap)
         }
      }
   }
   if len(lines) < 1 {
      return []map[string]interface{}{doc}, nil
   }
   return lines, nil
}

func init() {
   factory.RegisterFilter("split_tab", func(args []string) (lines api.Filter, err error) {
      filterSplitTab := &FilterSplitTab{}
      filterSplitTab.ParseOpts(args)
      return filterSplitTab, nil
   })
}

/////////////////////////////////////////////////
// split comma filter
/////////////////////////////////////////////////

type FilterSplitComma struct {
   BaseFilter
}

func (fs *FilterSplitComma) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   var lines []map[string]interface{}
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := strings.Split(docv.(string), ",")
         for idx := 0; idx < len(words); idx++ {
            newmap := util.Merge(make(map[string]interface{}), doc)
            newmap[fmt.Sprintf("%s.word", k)] = words[idx]
            lines = append(lines, newmap)
         }
      }
   }
   if len(lines) < 1 {
      return []map[string]interface{}{doc}, nil
   }
   return lines, nil
}

func init() {
   factory.RegisterFilter("split_comma", func(args []string) (lines api.Filter, err error) {
      filterSplitComma := &FilterSplitComma{}
      filterSplitComma.ParseOpts(args)
      return filterSplitComma, nil
   })
}

/////////////////////////////////////////////////
// split array filter
/////////////////////////////////////////////////

type FilterSplitArray struct {
   BaseFilter
}

func (fs *FilterSplitArray) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   var lines []map[string]interface{}
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         if val, ok := docv.([]string); ok {
            for idx := 0; idx < len(val); idx++ {
               lines = append(lines, map[string]interface{}{fmt.Sprintf("%s.item", k): val[idx]})
            }
         }
      }
   }
   if len(lines) < 1 {
      return []map[string]interface{}{doc}, nil
   }
   return lines, nil
}

func init() {
   factory.RegisterFilter("split_array", func(args []string) (lines api.Filter, err error) {
      filterSplitArray := &FilterSplitArray{}
      filterSplitArray.ParseOpts(args)
      return filterSplitArray, nil
   })
}

/////////////////////////////////////////////////
// field split line filter
/////////////////////////////////////////////////

type FilterFieldSplitLine struct {
   BaseFilter
}

func (fs *FilterFieldSplitLine) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := strings.Split(docv.(string), "\n")
         for idx := 0; idx < len(words); idx++ {
            doc[fmt.Sprintf("%s.f%d", k, idx+1)] = words[idx]
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_split_line", func(args []string) (lines api.Filter, err error) {
      filterFieldSplitLine := &FilterFieldSplitLine{}
      filterFieldSplitLine.ParseOpts(args)
      return filterFieldSplitLine, nil
   })
}

/////////////////////////////////////////////////
// field split word filter
/////////////////////////////////////////////////

type FilterFieldSplitWord struct {
   regex *regexp.Regexp
   BaseFilter
}

func (fs *FilterFieldSplitWord) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := fs.regex.Split(docv.(string), -1)
         for idx := 0; idx < len(words); idx++ {
            doc[fmt.Sprintf("%s.f%d", k, idx+1)] = words[idx]
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_split_word", func(args []string) (lines api.Filter, err error) {
      filterFieldSplitWord := &FilterFieldSplitWord{}
      filterFieldSplitWord.regex = regexp.MustCompile("\\W")
      filterFieldSplitWord.ParseOpts(args)
      return filterFieldSplitWord, nil
   })
}

/////////////////////////////////////////////////
// field split tab filter
/////////////////////////////////////////////////

type FilterFieldSplitTab struct {
   BaseFilter
}

func (fs *FilterFieldSplitTab) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := strings.Split(docv.(string), "\t")
         for idx := 0; idx < len(words); idx++ {
            doc[fmt.Sprintf("%s.f%d", k, idx+1)] = words[idx]
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_split_tab", func(args []string) (lines api.Filter, err error) {
      filterFieldSplitTab := &FilterFieldSplitTab{}
      filterFieldSplitTab.ParseOpts(args)
      return filterFieldSplitTab, nil
   })
}

/////////////////////////////////////////////////
// field split comma filter
/////////////////////////////////////////////////

type FilterFieldSplitComma struct {
   BaseFilter
}

func (fs *FilterFieldSplitComma) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         words := strings.Split(docv.(string), ",")
         for idx := 0; idx < len(words); idx++ {
            doc[fmt.Sprintf("%s.f%d", k, idx+1)] = words[idx]
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_split_comma", func(args []string) (lines api.Filter, err error) {
      filterFieldSplitComma := &FilterFieldSplitComma{}
      filterFieldSplitComma.ParseOpts(args)
      return filterFieldSplitComma, nil
   })
}

/////////////////////////////////////////////////
// field split array filter
/////////////////////////////////////////////////

type FilterFieldSplitArray struct {
   BaseFilter
}

func (fs *FilterFieldSplitArray) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         if val, ok := docv.([]string); ok {
            for idx := 0; idx < len(val); idx++ {
               doc[fmt.Sprintf("%s.f%d", k, idx+1)] = val[idx]
            }
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_split_array", func(args []string) (lines api.Filter, err error) {
      filterFieldSplitArray := &FilterFieldSplitArray{}
      filterFieldSplitArray.ParseOpts(args)
      return filterFieldSplitArray, nil
   })
}

/////////////////////////////////////////////////
// field array join comma filter
/////////////////////////////////////////////////

type FilterFieldArrayJoinComma struct {
   BaseFilter
}

func (fs *FilterFieldArrayJoinComma) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         if val, ok := docv.([]string); ok {
            doc[k] = strings.Join(val, ",")
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_array_join_comma", func(args []string) (lines api.Filter, err error) {
      filterFieldArrayJoinComma := &FilterFieldArrayJoinComma{}
      filterFieldArrayJoinComma.ParseOpts(args)
      return filterFieldArrayJoinComma, nil
   })
}

/////////////////////////////////////////////////
// field array join whitespace filter
/////////////////////////////////////////////////

type FilterFieldArrayJoinWhitespace struct {
   BaseFilter
}

func (fs *FilterFieldArrayJoinWhitespace) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         if val, ok := docv.([]string); ok {
            doc[k] = strings.Join(val, " ")
         }
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("field_array_join_whitespace", func(args []string) (lines api.Filter, err error) {
      filterFieldArrayJoinWhitespace := &FilterFieldArrayJoinWhitespace{}
      filterFieldArrayJoinWhitespace.ParseOpts(args)
      return filterFieldArrayJoinWhitespace, nil
   })
}

/////////////////////////////////////////////////
// digest filter
/////////////////////////////////////////////////

type FilterDigest struct {
   BaseFilter
}

func (fs *FilterDigest) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if docv, ok := doc[k]; ok && len(docv.(string)) > 0 {
         var hash []byte
         switch v {
         case "sha1":
            sha1hash := sha1.Sum([]byte(docv.(string)))
            hash = sha1hash[:]
            break
         case "md5":
            md5hash := md5.Sum([]byte(docv.(string)))
            hash = md5hash[:]
         default:
            panic(fmt.Sprintf("Unknown/unsupported hash func: %s", v))
         }
         doc[fmt.Sprintf("%s.md5", k)] = hex.EncodeToString(hash)
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("digest", func(args []string) (lines api.Filter, err error) {
      filterDigest := &FilterDigest{}
      filterDigest.ParseOpts(args)
      return filterDigest, nil
   })
}

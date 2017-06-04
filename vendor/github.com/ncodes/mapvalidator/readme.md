## MapValidator - A validator for maps

MapValidator is a simple package for validating map values. It allows a bunch of rules to be run through each map key/value checking whether rule is satisfied. All errors from the operation are returned. 

### Usage

```go
// Define rules
var rules = []Rule{
    TypeWithMsg("age", TypeMap, "data must be a map"),
}

// The map to validate
m := map[string]interface{}{
    "age": map[string]interface{}{},
}

// Validate does the work returning all errors
errs := Validate(m, rules...)
```

## Supported Rules

- RequiredWithMsg(fieldName, errorMsg)
- TypeWithMsg(fieldName, TypeChecker, message)
- EqualStringWithMsg(targetField, cmpVal, msg)
- EqualIntWithMsg(targetField, cmpVal, msg)
- EqualInt32WithMsg(targetField, cmpVal, msg)
- EqualInt64WithMsg(targetField, cmpVal, msg)
- EqualFloat32WithMsg(targetField, cmpVal, msg)
- EqualFloat64WithMsg(targetField, cmpVal, msg)
- Custom(targetField, msg, customFunc func(fieldValue interface{}, fullMap map[string]interface{}) bool)

## Supported TypeCheckers

- TypeString
- TypeNumber
- TypeDigit
- TypeMap
- TypeWithMsg

### TODO

Being a new package created to quickly solve a problem, it requires
many more new rules which will be added over time. Feel free to send PR for new
rules :)

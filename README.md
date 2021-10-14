# mqtt-rewriter

## Background

mqtt-rewriter is a tool that can forward data from a topic to another topic.

## Install

Todo...

## Usage

Currently only supports two types of forwarding:

- Delay
- Template

### Delay

Forward raw data from `$delay/interval/actual_topic` to `actual_topic` after `interval` milliseconds.

### Template

Use the go template to render the data or not and forward it to the corresponding topic.

#### Example

Use the following `config.yml`

```yml
rewriter:
  delay:
    enable: true
  template:
    enable: true
    rules:
      - topic: test
        type: raw
        targets:
          - topic: test2
      - topic: test2
        type: json
        targets:
          - topic: test3
            template: |-
              Get value: {{.msg}}
              Get nested value: {{.data.msg}}
              Get value of an array by index: {{index .array 0}}
              Get json of an object: {{json .data}}
              Get json of an array: {{json .array}}
              Range an array: 
              {{range $value := .array -}}
              {{println $value}}
              {{- end}}
              Range an array with index: 
              {{range $index,$value := .array -}}
              {{printf "%d:%s\n" $index $value}}
              {{- end}}
              Range an object:
              {{range $key,$value := .data -}}
              {{printf "%s:%s\n" $key $value}}
              {{- end}}
```

Send a data to `test`

```json
{
  "msg": "I'm msg.",
  "data": {
    "msg": "I'm msg of data.",
    "msg2": "I'm msg2 of data."
  },
  "array": ["I'm element of array[0]", "I'm element of array[1]"]
}
```

1. Match the first rule
   ```yml
   - topic: test
     type: raw
     targets:
       - topic: test2
   ```
   This rule type is raw, then program will forward raw data to topic `test2`.
2. Match the second rule when program forward raw data to `test2`

   ```yml
   - topic: test2
     type: json
     targets:
       - topic: test3
         template: |-
           Get value: {{.msg}}
           Get nested value: {{.data.msg}}
           Get value of an array by index: {{index .array 0}}
           Get json of an object: {{json .data}}
           Get json of an array: {{json .array}}
           Range an array: 
           {{range $value := .array -}}
           {{println $value}}
           {{- end}}
           Range an array with index: 
           {{range $index,$value := .array -}}
           {{printf "%d:%s\n" $index $value}}
           {{- end}}
           Range an object:
           {{range $key,$value := .data -}}
           {{printf "%s:%s\n" $key $value}}
           {{- end}}
   ```

   So we will received following data from topic `test3`

   ```
    Get value: I'm msg.
    Get nested value: I'm msg of data.
    Get value of an array by index: I'm element of array[0]
    Get json of an object: {"msg":"I'm msg of data.","msg2":"I'm msg2 of data."}
    Get json of an array: ["I'm element of array[0]","I'm element of array[1]"]
    Range an array:
    I'm element of array[0]
    I'm element of array[1]

    Range an array with index:
    0:I'm element of array[0]
    1:I'm element of array[1]

    Range an object:
    msg:I'm msg of data.
    msg2:I'm msg2 of data.
   ```

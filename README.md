# Multithreaded Roman Numeral to Integer Converter

Program to convert a Roman Numeral string to an integer using using channels and
concurrency.

The program works by spawning three threads with specific responsibilities:

- Read lexical values from a Roman Numeral string
- Look up numeric value for each lexeme
- Sum the numeric values

Data is fed from one Go routine to the next using independent channels.

`readSymbols` parses the input string creating tokens and sending them to
`readValues` over a channel. Within `readValues` each token is converted to a
numeric value using an internal map and sent to `addValues` using a channel.
`addValues` sums the values and sends the result to another channel. Convert
blocks until a value is received and returns that value.

Note: `readValues` creates a Go routine for each roman numeral that handles the
conversion to an integer. Thus the order in which values are summed in
`addValues` is asynchronous. Refer to the function itself or `TestReadSymbols`
for more details.

```text
+-------------------------+
| "MCMXCIV"             G |
|                         |
|                         |
| Convert()               |
+--+----------------------+
   |
   |        +--------------------------+        +--------------------------+
   |        | "MCMXCIV"              G |        | [M, CM, XC, ...]       G |
   |        |                          |        |                          |
   +------->|                          +------->|                          |
            | readSymbols()            |        | readValues()             |
            +--------------------------+        +-------------+------------+
                                                              |
                                                              v
      +-------------+------------+              +--------------------------+
      | [1000 + 900 + ...]     G |              | "CM"                   G |-+
      |                          |              |                          | |-+
      |                          |<-------------+                          | | |
      | addValues()              |              | anonymous()              | | |
      +--------------------------+              +--------------------------+ | |
                                                  +--------------------------+ |
                                                    +--------------------------+

G: Go Routine
```

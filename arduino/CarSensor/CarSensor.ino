#include <Arduino.h>

#include <WiFi.h>
#include <WiFiMulti.h>

#include <HTTPClient.h>
#include <NetworkClientSecure.h>

// https://wiki.seeedstudio.com/XIAO_ESP32C3_Getting_Started/
// https://wiki.seeedstudio.com/XIAO_ESP32C3_WiFi_Usage/
// https://stackoverflow.com/questions/65619359/how-to-add-headers-to-http-request-in-arduino-ide

// define led according to pin diagram in article
const int sleepEnable = D8;
const int led = D10; // there is no LED_BUILTIN available for the XIAO ESP32C3.

/*
// Create a new file called `secrets.h` and dump the following in it:
const char* ssid = "your_ssid";
const char* password = "your_wifi_password";
const char* webAddress = "https://yourWebsite.com/endpoint";
const String bearerToken = "your_bearerToken";
*/
#include "secrets.h"


// This is a GTS Root R1 cert, This certificate is valid until Sun, 22 Jun 2036 00:00:00 GMT
const char *rootCACertificate = "-----BEGIN CERTIFICATE-----\n"
"MIIFVzCCAz+gAwIBAgINAgPlk28xsBNJiGuiFzANBgkqhkiG9w0BAQwFADBHMQsw\n"
"CQYDVQQGEwJVUzEiMCAGA1UEChMZR29vZ2xlIFRydXN0IFNlcnZpY2VzIExMQzEU\n"
"MBIGA1UEAxMLR1RTIFJvb3QgUjEwHhcNMTYwNjIyMDAwMDAwWhcNMzYwNjIyMDAw\n"
"MDAwWjBHMQswCQYDVQQGEwJVUzEiMCAGA1UEChMZR29vZ2xlIFRydXN0IFNlcnZp\n"
"Y2VzIExMQzEUMBIGA1UEAxMLR1RTIFJvb3QgUjEwggIiMA0GCSqGSIb3DQEBAQUA\n"
"A4ICDwAwggIKAoICAQC2EQKLHuOhd5s73L+UPreVp0A8of2C+X0yBoJx9vaMf/vo\n"
"27xqLpeXo4xL+Sv2sfnOhB2x+cWX3u+58qPpvBKJXqeqUqv4IyfLpLGcY9vXmX7w\n"
"Cl7raKb0xlpHDU0QM+NOsROjyBhsS+z8CZDfnWQpJSMHobTSPS5g4M/SCYe7zUjw\n"
"TcLCeoiKu7rPWRnWr4+wB7CeMfGCwcDfLqZtbBkOtdh+JhpFAz2weaSUKK0Pfybl\n"
"qAj+lug8aJRT7oM6iCsVlgmy4HqMLnXWnOunVmSPlk9orj2XwoSPwLxAwAtcvfaH\n"
"szVsrBhQf4TgTM2S0yDpM7xSma8ytSmzJSq0SPly4cpk9+aCEI3oncKKiPo4Zor8\n"
"Y/kB+Xj9e1x3+naH+uzfsQ55lVe0vSbv1gHR6xYKu44LtcXFilWr06zqkUspzBmk\n"
"MiVOKvFlRNACzqrOSbTqn3yDsEB750Orp2yjj32JgfpMpf/VjsPOS+C12LOORc92\n"
"wO1AK/1TD7Cn1TsNsYqiA94xrcx36m97PtbfkSIS5r762DL8EGMUUXLeXdYWk70p\n"
"aDPvOmbsB4om3xPXV2V4J95eSRQAogB/mqghtqmxlbCluQ0WEdrHbEg8QOB+DVrN\n"
"VjzRlwW5y0vtOUucxD/SVRNuJLDWcfr0wbrM7Rv1/oFB2ACYPTrIrnqYNxgFlQID\n"
"AQABo0IwQDAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4E\n"
"FgQU5K8rJnEaK0gnhS9SZizv8IkTcT4wDQYJKoZIhvcNAQEMBQADggIBAJ+qQibb\n"
"C5u+/x6Wki4+omVKapi6Ist9wTrYggoGxval3sBOh2Z5ofmmWJyq+bXmYOfg6LEe\n"
"QkEzCzc9zolwFcq1JKjPa7XSQCGYzyI0zzvFIoTgxQ6KfF2I5DUkzps+GlQebtuy\n"
"h6f88/qBVRRiClmpIgUxPoLW7ttXNLwzldMXG+gnoot7TiYaelpkttGsN/H9oPM4\n"
"7HLwEXWdyzRSjeZ2axfG34arJ45JK3VmgRAhpuo+9K4l/3wV3s6MJT/KYnAK9y8J\n"
"ZgfIPxz88NtFMN9iiMG1D53Dn0reWVlHxYciNuaCp+0KueIHoI17eko8cdLiA6Ef\n"
"MgfdG+RCzgwARWGAtQsgWSl4vflVy2PFPEz0tv/bal8xa5meLMFrUKTX5hgUvYU/\n"
"Z6tGn6D/Qqc6f1zLXbBwHSs09dR2CQzreExZBfMzQsNhFRAbd03OIozUhfJFfbdT\n"
"6u9AWpQKXCBfTkBdYiJ23//OYb2MI3jSNwLgjt7RETeJ9r/tSQdirpLsQBqvFAnZ\n"
"0E6yove+7u7Y/9waLd64NnHi/Hm3lCXRSHNboTXns5lndcEZOitHTtNCjv0xyBZm\n"
"2tIMPNuzjsmhDYAPexZ3FL//2wmUspO8IFgV6dtxQ/PeEMMA3KgqlbbC1j+Qa3bb\n"
"bP6MvPJwNQzcmRk13NfIRmPVNnGuV/u3gm3c\n"
"-----END CERTIFICATE-----\n";



void setup() {
    // initialize digital pin led as an output
  pinMode(led, OUTPUT);
  pinMode(sleepEnable, INPUT);
  pinMode(A0, INPUT);         // ADC
  digitalWrite(led, LOW);    // turn the LED off

  Serial.begin(115200);
  delay(100);

  // We start by connecting to a WiFi network
  Serial.println();
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  int delaySum = 0;
  while (WiFi.status() != WL_CONNECTED) {
    delaySum += 500;
    delay(500);
    Serial.print(".");

    if ( delaySum > 10000) {
      break;
    }
  }

  if (WiFi.status() == WL_CONNECTED) {
    delay(1000);
    digitalWrite(led, HIGH);   // turn the LED on 
    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
  }

  
}

void loop() {
  uint32_t Vbatt = 0;
  int measurementCount = 32;
  for(int i = 0; i < measurementCount; i++) {
    Vbatt = Vbatt + analogReadMilliVolts(A0); // ADC with correction
  }

  // 220k Ohm and 1M Ohm
  float Vbattf = 5.5454545 * float(Vbatt) / float(measurementCount) / 1000.0;     // attenuation ratio, mV --> V
  Serial.println(Vbattf, 3);
  // delay(1000);


  NetworkClientSecure *client = new NetworkClientSecure;
  if (client && WiFi.status() == WL_CONNECTED) {
    client->setCACert(rootCACertificate);
    // client->setInsecure();
    {
      // Add a scoping block for HTTPClient https to make sure it is destroyed before NetworkClientSecure *client is
      HTTPClient https; //Object of class HTTPClient
      if(https.begin(*client, webAddress)){
        https.addHeader("Authorization", "Bearer " + bearerToken);

        // {"measurements":[{"valName":"battV","valFloat":1.1},{"valName":"temp","valFloat":34.1}]}
        String sendBody = "{\"measurements\":[";
        
        sendBody += "{\"valName\": \"BattV\", \"valFloat\": ";
        sendBody += String(Vbattf, 3);
        sendBody += "}";

        sendBody += "]}";
        Serial.println(sendBody);

        int httpCode = https.POST(sendBody);
        String returnBody = https.getString();
        Serial.println(httpCode);
        if (httpCode > 0) {
          Serial.printf("[HTTPS] POST... code: %d, body: %s\n", httpCode, returnBody);
        } else {
          Serial.printf("[HTTPS] POST... failed, error: %d, %s\n", httpCode, https.errorToString(httpCode).c_str());
        }
        https.end(); //Close connection
      }else {
        Serial.printf("[HTTPS] Unable to connect\n");
      }

    }

    delete client;
  } else {
    Serial.println("Unable to create client");
  }




  while(!digitalRead(sleepEnable)) {
    delay(100);
  }

  //Go to sleep now
  digitalWrite(led, LOW);    // turn the LED off
  esp_sleep_enable_timer_wakeup(10*60000000); // 10 minutes
  Serial.println("Going to sleep now");
  esp_deep_sleep_start();
}

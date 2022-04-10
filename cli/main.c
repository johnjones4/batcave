#include "main.h"
#include <stdio.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#define PORT 8090


int main(int argc, char const *argv[])
{
  int sock = 0, valread;
  struct sockaddr_in serv_addr;
  if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0)
  {
      printf("\n Socket creation error \n");
      return -1;
  }

  serv_addr.sin_family = AF_INET;
  serv_addr.sin_port = htons(PORT);
      
  if(inet_pton(AF_INET, "127.0.0.1", &serv_addr.sin_addr)<=0) 
  {
    return -1;
  }
  
  if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0)
  {
    return -1;
  }

  if (login(sock) < 0) {
    return -1;
  }

  return runloop(sock);
}

int login(int sock) {
  char *username = "John";
  char *password = "banana";
  char *payload = malloc(sizeof(username) + sizeof(password) + 2);
  sprintf(payload, "%s:%s\n", username, password);  
  send(sock, payload, strlen(payload), 0);
  free(payload);
  char *buffer = malloc(sizeof(char) * 4096);
  readNext(sock, buffer);
  int c = strcmp(buffer, "ok");
  free(buffer);
  return c;
}

int readNext(int sock, char *buffer) {
  buffer[0] = 0;
  char lastChar = ' ';
  int ptr = 0;
  while (lastChar != '\r') {
    char buf[1024] = {0};
    int n = read(sock, buf, 256);
    lastChar = buf[n - 1];
    strcpy(buffer, buf);
    ptr += n;
  }
  buffer[ptr - 1] = 0;
  return ptr - 1;
}

int runloop(int sock) {
  while (1) {
    char inputBuffer[1024];
    printf("You> ");
    scanf("%s",&inputBuffer);

    if (strcmp(inputBuffer, "quit") == 0) {
      return 0;
    }

    int l = strlen(inputBuffer);
    inputBuffer[l] = '\n';
    inputBuffer[l + 1] = 0;

    send(sock, inputBuffer, strlen(inputBuffer), 0);

    char *buffer = malloc(sizeof(char) * 4096);
    readNext(sock, buffer);
    free(buffer);

    printf("HAL>%s\n", buffer);
  }
}
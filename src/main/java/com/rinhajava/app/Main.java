package com.rinhajava.app;


import com.rinhajava.app.http.Server;

import java.io.IOException;

public class Main {
    public static void main(String[] args ) throws IOException {
        Server server = new Server();
        server.startServerHandlers();
        System.out.println("server started");
    }
}

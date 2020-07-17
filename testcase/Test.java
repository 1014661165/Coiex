package com;

import java.io.IOException;

public class Test extends A implements B, C{

    public String a = "public int \"add(int a, int b){return a+b;}";

    public static void main(String[] args){
        String a = "public static int add(int a,int b){return a+b;}";
    }

    public int add(List<Integer> list){
        return 0;
    }
}
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
class Login extends StatefulWidget {
  @override
  _LoginState createState() => _LoginState();
}

class _LoginState extends State<Login> {
  @override

  Color backgroundCol = new Color(0xff222831);
  Color appBarCol = new Color(0xff393e46);
  Color primary = new Color(0xff00adb5);
  Color tiles = new Color(0xff3e4149);
  Color secondary = new Color(0xffeeeeee);
  Color listText = new Color(0xffdbe2ef);
  
  Widget build(BuildContext context) {
    var screenWidth = MediaQuery.of(context).size.width;
    var screenHeight = MediaQuery.of(context).size.height;
    return Scaffold(
      backgroundColor: backgroundCol,
      body: Container(
        child: Center(
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Center(
                    child: Image(image: AssetImage('assets/listIcon.png'),
                    width: screenWidth/9,)
                  ),
                  Container(width: 30,),
                  Text("MyTodo", style: GoogleFonts.openSans(color: Colors.white, fontSize: screenHeight/20),)
                ],
              ),
              Container(
                height: screenHeight/1.5,
                width: 1,
                color: Colors.grey[600],
              ),
              Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    margin: EdgeInsets.symmetric(vertical: 60),
                    child: Text("Login", style: GoogleFonts.openSans(fontSize: 40, color: Colors.white),)
                  ),
                  Container(
                    height: screenHeight/15,
                    width: screenWidth/5,
                    child: TextField(
                      decoration: InputDecoration(
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.all(Radius.circular(30))
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.all(Radius.circular(30)),
                          borderSide: BorderSide(
                            color: Colors.white
                          )
                        ),
                        filled: true,
                        fillColor: appBarCol,
                        hintText: "Username",
                        hintStyle: GoogleFonts.openSans(color: Colors.white),
                      ),
                    ),
                  ),
                  Container(
                    height: 0,
                  ),
                  Container(
                    height: screenHeight/15,
                    width: screenWidth/5,
                    child: TextField(
                      decoration: InputDecoration(
                        focusColor: Colors.white,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.all(Radius.circular(30))
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.all(Radius.circular(30)),
                          borderSide: BorderSide(
                            color: Colors.white
                          )
                          
                        ),
                        hintText: "Password",
                        hintStyle: GoogleFonts.openSans(color: Colors.white),
                        filled: true,
                        fillColor: appBarCol,
                      ),
                    ),
                  )
                ],
              ),
            ],
          ),
        ),
      ),
      
    );
  }
}
Test setup.

  $ FW="$TESTDIR/../farnsworth"

Create an archive of the Java project.

  $ "$FW" archive --project "$TESTDIR/projects/java" --public java.zip
  $ ls .
  java.zip

Decompress the archive.

  $ unzip -qq java.zip
  $ ls .
  java
  java.zip

Check that ignored files were ignored

  $ ls java
  README.md
  build.gradle
  src

Check the Main.java file.

  $ cat java/src/main/java/Main.java
  public class Main {
      public static void main(String[] args) {
      }
  }

Check the MainTest.java file.

  $ cat java/src/test/java/MainTest.java
  import org.junit.Test;
  import static org.junit.Assert.*;
  
  public class MainTest {
      @Test public void dummyTestPublic() {
          assertNotNull("dummy test public", 1 + 2);
      }
  }


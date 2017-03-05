import org.junit.Test;
import static org.junit.Assert.*;

public class MainTest {
    @Test public void dummyTestPublic() {
        assertNotNull("dummy test public", 1 + 2);
    }
    //++ hide

    @Test public void dummyTestPrivate() {
        assertNotNull("dummy test private", 2 + 3);
    }
    //++ stop
}


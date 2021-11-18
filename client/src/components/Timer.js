import {CircularProgressbar, buildStyles} from "react-circular-progressbar"
import "react-circular-progressbar/dist/styles.css"

const purple="#645394";
const red = '#f54e4e';

function Timer(){
    //const[timerMinutes, setTimerMinutes] = useState('00');
    //const [timerSeconds, settimerSeconds] = useState('00');
    //let interval= userRef();
    
    const percentage = 100; 
    const minutes = 1 
    const seconds = '00'

    return(
        <div>
            <CircularProgressbar
                value={percentage}
                text={minutes + ':' + seconds}
                styles={buildStyles({
                textColor:'#fff',
                tailColor:'rgba(255,255,255,.2)',})} 
        />
        </div>
    )
}
export default Timer;
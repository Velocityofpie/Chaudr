import {CircularProgressbar, buildStyles} from "react-circular-progressbar"
import "react-circular-progressbar/dist/styles.css"

const purple="#645394";
const red = '#f54e4e';

function Timer(){
    return(
        <div>
            <CircularProgressbar value={60} text={`${60}`} 
            styles={buildStyles({
                textColor:'#fff',
                
                tailColor:'rgba(255,255,255,.2)',
                })} />
        </div>
    )
}
export default Timer;
import React from 'react'

function Popup(props) {
    return (props.trigger)?(
        <div className="popup">
           <div className="popup__inner">
                <button className="close_btn">close</button>
                {props.children}
           </div>
        </div>
    ):"";
}

export default Popup

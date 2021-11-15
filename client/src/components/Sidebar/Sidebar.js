import React from 'react'
import "./Sidebar.css";
import AddIcon from "@material-ui/icons/Add";
import SidebarChannel from './SidebarChannel';
import Settings from '@material-ui/icons/Settings';
import {Avatar} from '@material-ui/core/';
import {useState} from 'react'



function Sidebar() {
    const[buttonPopup,settButtonPopup]=useState(false)
    const[buttonPopup2,settButtonPopup2]=useState(false)
    const[buttonPopup3,settButtonPopup3]=useState(false)
    return (
        <div className="sidebar">
            <div className="sidebar__top">
                <h3>Chauder</h3>
            </div>
            <div className="sidebar__channels">
                <div className="sidebar__channelsHeader">
                    <div className="sidebar__header">      
                        <h3>Chat rooms</h3>
                    </div>
                    <button onClick={()=> settButtonPopup(true)} className="sidebar__addChannel" icon={<AddIcon/>} >
                    </button>
                    <Popup trigger={buttonPopup} setTrigger={settButtonPopup}>
                        <button onClick={()=> settButtonPopup2(true)} className="Join_btn">
                            Join New Chat
                        </button >
                        <button onClick={()=> settButtonPopup3(true)} className="Make_btn">
                            Make New Chat
                        </button>
                    </Popup>
                    <Join_Popup trigger={buttonPopup2} setTrigger={settButtonPopup2} >
                    </Join_Popup>
                    <Make_Popup trigger={buttonPopup3} setTrigger={settButtonPopup3} >
                    </Make_Popup>
                </div>
                    <div className="sidebar__channelsList">
                        <SidebarChannel />
                        <SidebarChannel />
                        <SidebarChannel />
                    </div>
                </div>
                  
                <div className="sidebar__bottom"> 
                    <button>Logout</button>
                    <Avatar src={"/chaudrlogo.png"} style={{ height: '75px', width: '75px' } }></Avatar>
                           
                </div>  
            </div>
            
    )
}


function Popup(props) {
    return (props.trigger)?(
        <div className="popup">
           <div className="popup__inner">
                <button className="close_butn" onClick={()=>props.setTrigger(false)}>close</button> 
                {props.children}
           </div>
        </div>
    ):"";
}

function Join_Popup(props) {
    return (props.trigger)?(
        <div className="join_popup">
           <div className="join_popup__inner">
                <button className="close_butn2" onClick={()=>props.setTrigger(false)}>Back</button>
                    Join New Chat
                {props.children}
           </div>
        </div>
    ):"";
}
function Make_Popup(props) {
    return (props.trigger)?(
        <div className="make_popup">
           <div className="make_popup__inner">
                <button className="close_butn3" onClick={()=>props.setTrigger(false)}>Back</button>
                <input 
                type="text"
                placeholder={'Enter Channel Name'} 
                label="ChannelName"  
                margin="normal"/>
                {props.children}
                <button>Create</button>
           </div>
        </div>
    ):"";
}

export default Sidebar

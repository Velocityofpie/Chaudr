import React from 'react'
import "./Sidebar.css";
import AddIcon from "@material-ui/icons/Add";
import SidebarChannel from './SidebarChannel';
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
                    <button onClick={()=> settButtonPopup(true)} className="sidebar__addChannel"  >
                        <span className="button__icon">
                                <AddIcon/>
                        </span>
                    </button>
                    <Popup trigger={buttonPopup} setTrigger={settButtonPopup} icon={<AddIcon/>}>
                        <h3>Create or Join</h3>
                        <button onClick={()=> settButtonPopup2(true)} className="Join_btn">
                            Join New Chat
                        </button >
                        <button onClick={()=> settButtonPopup3(true)} className="Create_btn">
                            Create New Chat
                        </button>
                    </Popup>
                    <Join_Popup trigger={buttonPopup2} setTrigger={settButtonPopup2} >
                    </Join_Popup>
                    <Create_Popup trigger={buttonPopup3} setTrigger={settButtonPopup3} >
                    </Create_Popup>
                </div>
                    <div className="sidebar__channelsList">
                        <SidebarChannel />
                        <SidebarChannel />
                        <SidebarChannel />
                    </div>
                </div>
                  
                <div className="sidebar__bottom"> 
                    <button>Logout</button>
                    <Avatar src={"/chaudrlogo.png"} style={{ height: '70px', width: '70px' } }></Avatar> 
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
                <h2>Room ID</h2>
                <input 
                type="text"
                placeholder={'Enter Room ID'} 
                label="ChannelName"  
                margin="normal"/>
                <h2>Code</h2>
                <input 
                type="text"
                placeholder={'Enter Code'} 
                label="ChannelName"  
                margin="normal"/>
                <h2>Channel Name</h2>
                <input 
                type="text"
                placeholder={'Enter Channel Name'} 
                label="ChannelName"  
                margin="normal"/>
                <button>join</button>
                {props.children}
           </div>
        </div>
    ):"";
}
function Create_Popup(props) {
    return (props.trigger)?(
        <div className="create_popup">
           <div className="create_popup__inner">
                <button className="close_butn3" onClick={()=>props.setTrigger(false)}>Back</button>
                <h2>Channel Name</h2>
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

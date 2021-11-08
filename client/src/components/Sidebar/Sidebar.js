import React from 'react'
import "./Sidebar.css";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import AddIcon from "@material-ui/icons/Add";
import SidebarChannel from './SidebarChannel';
import Settings from '@material-ui/icons/Settings';


function Sidebar() {
    return (
        <div className="sidebar">
            <div className="sidebar__top">
                <h3>Chauder</h3>  
                <Settings />   
            </div>
            <div className="sidebar__channels">
                <div className="sidebar__channelsHeader">
                    <div className="sidebar__header">      
                        <h4>Chat rooms</h4>
                    </div>
                    <AddIcon className="sidebar__addChannel" />
                </div>
                    <div className="sidebar__channelsList">
                        <SidebarChannel />
                        <SidebarChannel />
                        <SidebarChannel />
                    </div>
                </div>
                <div className="sidebar__bottom">
                    <h3>bottom half</h3>  
                        <Settings />   
                </div>  
            </div>
            
    )
}

export default Sidebar

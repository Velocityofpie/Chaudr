import React,{useState}  from 'react'
import './ChatHeader.css'
import ExpandMoreIcon from "@material-ui/icons/ExpandMore"



function ChatHeader() {
    return (
        <div className='chatHeader'>
            <div className='chatHeader__left'>
                <h3>
                    <span className='chatHeader_name'> 
                        -
                    </span>
                        test Channel Name
                </h3>
            </div>
            <Navbar>
                <NavItem icon={<ExpandMoreIcon/>}>
                    <DropdownMenu/>
                </NavItem>
            </Navbar>
            
        </div>
    )
}
function DropdownMenu(){
    function DropdownItem(props){
        return(
            <a href="#" className="menu_item" >
                
                {props.children}
            </a>
        )
    }
    return(
        <div className="dropdown">

        </div>
    )
}


function Navbar(props) {
    return (
        <nav className="navbar">
            <ul className="navbar_nav"> {props.children }</ul>
        </nav>
    )
}

function NavItem(props) {
    const[open,setOpen]=useState(false);
    return (
        <li className="nav_item">
            <a href="#" className="icon_button" onClick={()=> setOpen(!open)}>
                {props.icon}
            </a>
            {open && props.children}
        </li>
    )
}

export default ChatHeader


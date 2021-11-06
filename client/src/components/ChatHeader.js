import React from 'react'
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
            <div className='chatHeader__right'>
                Chat setting
                < ExpandMoreIcon/>
                
            </div>

        </div>
    )
}

export default ChatHeader

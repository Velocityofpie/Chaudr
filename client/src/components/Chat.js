import React from 'react'
import './Chat.css'
import ChatHeader from './ChatHeader'
import EmojiEmotionsIcon from '@material-ui/icons/EmojiEmotions'

function Chat() {
    return (
        <div className='chat'>
            <ChatHeader/>
            <div className='chat__messages'>

            </div>
            <div className='chat__input'>

                <form>
                    <input placeholder={'input text'}/>
                    <button className="chat__inputButton" type ="submit">
                        send message
                    </button>
                </form>
                <div className="chat__inputIcons">
                    <EmojiEmotionsIcon frontsize="large"  />
                    
                </div>
            </div>
        </div>
    )
}

export default Chat

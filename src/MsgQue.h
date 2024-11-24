#ifndef MSG_QUE_H
#define MSG_QUE_H

#include <future>

#include "message/Msg.h"

class MsgRegister {

   std::shared_ptr<EuloMsgType> msg;
public:

   static MsgRegister & obj() {
      static MsgRegister ins;
      return ins;
   }

   void setMsg (std::shared_ptr<EuloMsgType> msg){
      this->msg = msg;
   }

   std::shared_ptr<EuloMsgType> getMsg(){
      return this->msg;
   }

};


#endif
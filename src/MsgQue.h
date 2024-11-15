#ifndef MSG_QUE_H
#define MSG_QUE_H

#include <future>

#include "Msg.h"

class MsgRegister {

   std::shared_ptr<NormalMsg> msg;
public:

   static MsgRegister & obj() {
      static MsgRegister ins;
      return ins;
   }

   void setMsg (std::shared_ptr<NormalMsg> msg){
      this->msg = msg;
   }

   std::shared_ptr<NormalMsg> getMsg(){
      return this->msg;
   }

};


#endif
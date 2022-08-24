import type { MessageApiInjection } from 'naive-ui/lib/message/src/MessageProvider';
import type { DialogApiInjection } from 'naive-ui/lib/dialog/src/DialogProvider';

declare global {
  interface Window {
    $message: MessageApiInjection;
    $dialog: DialogApiInjection;
  }
}

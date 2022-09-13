// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import * as mqtt from 'mqtt/dist/mqtt.min';
import type { ISubscriptionGrant } from 'mqtt';
import { ref } from 'vue';

export const logs = ref<string[]>([]);

export const useMqtt = () => {
  const client = mqtt.connect('ws://127.0.0.1:19039', {
    // clean: false,
    clientId: localStorage.getItem('mqtt-client-id') || undefined,
  });

  localStorage.setItem('mqtt-client-id', client.options.clientId);

  client.on('connect', () => {
    // eslint-disable-next-line no-console
    console.log('Connected to MQTT');

    client.subscribe(`log/${client.options.clientId}`, { qos: 2 }, (err: Error, granted: ISubscriptionGrant[]) => {
      if (err) {
        // eslint-disable-next-line no-console
        console.error(err);
      } else {
        // eslint-disable-next-line no-console
        console.log('Subscribed to topic', granted);
      }
    });
  });

  client.on('error', (error: Error) => {
    // eslint-disable-next-line no-console
    console.log('MQTT error', error);
  });

  client.on('message', (topic: string, payload: ArrayBuffer) => {
    logs.value.push(payload.toString());
    // console.log('MQTT message', topic, payload.toString(), packet);
  });

  return client;
};

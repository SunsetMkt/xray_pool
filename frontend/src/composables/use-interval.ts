import { onBeforeUnmount, ref } from 'vue';

const useInterval = (fn: any, ms: number, autoStart = true) => {
  const timer = ref<number | null>(null);
  if (autoStart) {
    timer.value = setInterval(() => {
      fn();
    }, ms);
    fn();
  }
  const resetInterval = () => {
    if (timer.value === null) return;
    clearInterval(timer.value);
    timer.value = setInterval(() => {
      fn();
    }, ms);
    fn();
  };
  const stopInterval = () => {
    if (timer.value === null) return;
    clearInterval(timer.value);
  };
  onBeforeUnmount(() => {
    if (timer.value === null) return;
    clearInterval(timer.value);
  });
  return {
    timer,
    resetInterval,
    stopInterval,
  };
};

export default useInterval;

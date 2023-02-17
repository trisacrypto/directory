import React from "react";

export type UseScriptStatus = 'loading' | 'idle' | 'ready' | 'error'

const useScript = (src: string) => {
  const [status, setStatus] = React.useState<UseScriptStatus>(src ? 'loading' : 'idle');

  React.useEffect(() => {
    if (!src) {
      setStatus('idle');
      return;
    }

    let script: any = document.querySelector(`script[src="${src}"]`);

    if (!script) {
      script = document.createElement('script');
      script.src = src;
      script.async = true;
      script.setAttribute('data-status', 'loading');
      document.body.appendChild(script);

      const setDataStatus = (event: any) => {
        script.setAttribute(
          'data-status',
          event.type === 'load' ? 'ready' : 'error'
        );
      };
      script.addEventListener('load', setDataStatus);
      script.addEventListener('error', setDataStatus);
    } else {
      setStatus(script.getAttribute('data-status'));
    }

    const setStateStatus = (event: any) => {
      setStatus(event.type === 'load' ? 'ready' : 'error');
    };

    script.addEventListener('load', setStateStatus);
    script.addEventListener('error', setStateStatus);

    return () => {
      if (script) {
        script.removeEventListener('load', setStateStatus);
        script.removeEventListener('error', setStateStatus);
      }
    };
  }, [src]);

  return status;
};

export default useScript
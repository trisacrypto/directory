import React from "react";

function useSessionStorageState(key, defaultValue) {
    const [state, setState] = React.useState(() => {
        let value
        try {
            value = JSON.parse(
                window.sessionStorage.getItem(key) || JSON.stringify(defaultValue),
            )
        } catch (e) {
            value = defaultValue
        }
        return value
    })

    React.useEffect(
        () => {
            window.localStorage.setItem(key, state)
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [state],
    )
    return [state, setState]
}

export default useSessionStorageState;

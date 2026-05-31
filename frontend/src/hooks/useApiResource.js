import React from "react";

export function useApiResource(loader) {
  const [state, setState] = React.useState({
    data: null,
    error: "",
    loading: true,
  });

  React.useEffect(() => {
    let ignore = false;

    async function load() {
      try {
        const data = await loader();
        if (!ignore) {
          setState({ data, error: "", loading: false });
        }
      } catch (error) {
        if (!ignore) {
          setState({ data: null, error: error.message, loading: false });
        }
      }
    }

    load();

    return () => {
      ignore = true;
    };
  }, [loader]);

  return state;
}

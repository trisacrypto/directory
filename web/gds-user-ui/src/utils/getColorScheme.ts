const getColorScheme = (status: string | boolean) => {
  if (status === 'yes' || status === true) {
    return 'green';
  } else {
    return 'red';
  }
};

export default getColorScheme;

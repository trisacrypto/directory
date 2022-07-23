const getColorScheme = (status: string | boolean) => {
  if (status === 'yes' || status === true) {
    return 'green';
  } else {
    return 'orange';
  }
};

export default getColorScheme;

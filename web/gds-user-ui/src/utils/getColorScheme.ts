const getColorScheme = (status: string | boolean) => {
  console.log('status', status);
  if (status === 'yes' || status === true) {
    return 'green';
  } else {
    return 'red';
  }
};

export default getColorScheme;

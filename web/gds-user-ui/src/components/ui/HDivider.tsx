import { Text, TextProps } from '@chakra-ui/react';

const HDivider = (props: TextProps) => {
  return <Text {...props}>{' | '}</Text>;
};

export default HDivider;

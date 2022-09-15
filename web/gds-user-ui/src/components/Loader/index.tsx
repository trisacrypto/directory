import { Stack, Spinner, Flex, Box, Text, VStack, useColorModeValue } from '@chakra-ui/react';

interface LoaderProps {
  text?: string;
}
const Loader: React.FC<LoaderProps> = ({ text, ...rest }) => {
  return (
    <Flex height={'100vh'} alignItems={'center'} textAlign={'center'} justifyContent={'center'}>
      <VStack spacing={4}>
        <Spinner color={useColorModeValue('blue.500', 'whiteAlpha.500')} size="xl" {...rest} />
        <Text>{text}</Text>
      </VStack>
    </Flex>
  );
};
Loader.defaultProps = {
  text: 'Loading...'
};

export default Loader;

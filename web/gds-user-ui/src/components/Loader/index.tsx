import { Stack, Spinner, Flex, Box, Text, VStack } from '@chakra-ui/react';

interface LoaderProps {
  text?: string;
}
const Loader: React.FC<LoaderProps> = (props) => {
  return (
    <Flex
      height={'100vh'}
      bg={'white'}
      alignItems={'center'}
      textAlign={'center'}
      justifyContent={'center'}>
      <VStack spacing={4}>
        <Spinner color="blue.500" size="xl" />
        <Text>{props?.text}</Text>
      </VStack>
    </Flex>
  );
};
Loader.defaultProps = {
  text: 'Loading...'
};

export default Loader;

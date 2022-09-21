import { Stack, Spinner, Flex, Box, Text, VStack } from '@chakra-ui/react';

interface LoaderProps {
  text?: string;
}
const Loader: React.FC<LoaderProps> = ({ text, ...rest }) => {
  return (
    <Flex height={'100vh'} alignItems={'center'} textAlign={'center'} justifyContent={'center'}>
      <VStack spacing={4}>
        <Spinner color="blue.500" size="xl" {...rest} />
        <Text>{text}</Text>
      </VStack>
    </Flex>
  );
};
Loader.defaultProps = {
  text: 'Loading...'
};

export default Loader;

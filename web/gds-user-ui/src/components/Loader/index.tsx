import { Spinner, Flex, Text, VStack } from '@chakra-ui/react';

interface LoaderProps {
  text?: string;
  h?: string;
  withoutText?: boolean;
}
const Loader: React.FC<LoaderProps> = ({ text, h, withoutText = false, ...rest }) => {
  return (
    <Flex
      height={h || '100vh'}
      alignItems={'center'}
      textAlign={'center'}
      justifyContent={'center'}>
      <VStack spacing={4}>
        <Spinner color="blue.500" size="xl" {...rest} />
        {withoutText ? null : <Text>{text}</Text>}
      </VStack>
    </Flex>
  );
};
Loader.defaultProps = {
  text: 'Loading...'
};

export default Loader;

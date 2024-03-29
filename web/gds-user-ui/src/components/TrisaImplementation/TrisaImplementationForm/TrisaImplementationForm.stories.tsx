import { Meta, Story } from '@storybook/react';
import TrisaImplementationForm from '.';

type TrisaImplementationFormProps = {
  headerText: string;
  name: string;
  type: 'TestNet' | 'MainNet';
};

export default {
  title: 'components/Trisa Implementation Form',
  component: TrisaImplementationForm
} as Meta;

const Template: Story<TrisaImplementationFormProps> = (args:any) => (
  <TrisaImplementationForm {...args} />
);

export const Default = Template.bind({});
Default.args = {
  headerText: 'TRISA Endpoint: MainNet'
};

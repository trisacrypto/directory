import { Meta, Story } from '@storybook/react';
import DeleteButton, { DeleteButtonProps } from './DeleteButton';

export default {
  title: 'components/DeleteButton',
  component: DeleteButton
} as Meta;

const Template: Story<DeleteButtonProps> = (args) => <DeleteButton {...args} />;

export const Default = Template.bind({});
Default.args = {};

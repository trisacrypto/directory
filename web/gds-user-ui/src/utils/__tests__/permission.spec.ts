import { hasPermission } from '../permission';


jest.mock('application/store', () => ({
    getState: jest.fn().mockReturnValue({
        user: {
            user: {
                permission: ['read:collaborators', 'create:collaborators', 'update:collaborators', 'approve:collaborators', 'read:certificates', 'create:certificates', 'update:certificates', 'read:vasp', 'create:vasp', 'update:vasp']
            }
        }
    })
}));
describe('permission func handler', () => {
    beforeAll(() => {

    });

    it('should return true if user has read collaborotor permission', () => {
        expect(hasPermission('read:collaborators')).toBeTruthy();
    });

    it('should return false if user does not have permission', () => {
        expect(hasPermission('revoke:certificates')).toBeFalsy();
    });
});

package leveldb

func (s *leveldbTestSuite) TestObjectIterator() {
	iter := s.db.Iter("")
	defer iter.Release()

	objs := 0
	for iter.Next() {
		objs++
		key := iter.Key()
		obj, err := iter.Object(false)
		s.NoError(err)
		s.Equal(key, obj.Key)
		s.Empty(obj.Data)
	}

	s.NoError(iter.Error())
	s.Equal(2, objs)
}

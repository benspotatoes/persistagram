package backend

func (b *backendImpl) Exists(data InstagramMetadata) bool {
	_, err := b.db.Stat(data.remoteFilename())
	return err == nil
}

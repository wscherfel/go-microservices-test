package model

import(
  "github.com/jinzhu/gorm"
  "time"
  )


// StringModelDAO is a data access object to a database containing StringModels
type StringModelDAO struct {
  db *gorm.DB
}

// NewStringModelDAO creates a new Data Access Object for the
// StringModel model.
func NewStringModelDAO (db *gorm.DB) *StringModelDAO {
  return &StringModelDAO{
    db:db,
  }
}

// Create will create single StringModel in database.
func (dao *StringModelDAO) Create(m *StringModel) {
  dao.db.Create(m)
}

// Read will find all DB records matching
// values in a model given by parameter
func (dao *StringModelDAO) Read(m *StringModel) []StringModel {
  retVal := []StringModel{}
  dao.db.Where(m).Find(&retVal)

  return retVal
}

// ReadByID will find StringModel by ID given by parameter
func (dao *StringModelDAO) ReadByID(id uint) *StringModel{
  m := &StringModel{}
  if dao.db.First(&m, id).RecordNotFound() {
    return nil
  }

  return m
}

// Update will update a record of StringModel in DB
func (dao *StringModelDAO) Update(m *StringModel, id uint) *StringModel{
  oldVal := dao.ReadByID(id)
  if oldVal == nil {
    return nil
  }

  dao.db.Model(&oldVal).Updates(m)
  return oldVal
}

// UpdateAllFields will update ALL fields of StringModel in db
// with values given in the StringModel by parameter
func (dao *StringModelDAO) UpdateAllFields(m *StringModel) *StringModel{
	dao.db.Save(&m)
	return m
}

// Delete will soft-delete a single StringModel
func (dao *StringModelDAO) Delete(m *StringModel) {
  dao.db.Delete(m)
}

// GetUpdatedAfter will return all StringModels that were
// updated after given timestamp
func (dao *StringModelDAO) GetUpdatedAfter(timestamp time.Time) []StringModel {
	m := []StringModel{}
	dao.db.Where("updated_at > ?", timestamp).Find(&m)
	return m
}

// GetAll will return all records of StringModel in database
func (dao *StringModelDAO) GetAll() []StringModel {
	m := []StringModel{}
	dao.db.Find(&m)

	return m
}
// ReadByKey will find all records
// matching the value given by parameter
func (dao *StringModelDAO) ReadByKey (m string) []StringModel {
  retVal := []StringModel{}
  dao.db.Where(&StringModel{ Key : m }).Find(&retVal)

  return retVal
}

// DeleteByKey deletes all records in database with
// Key the same as parameter given
func (dao *StringModelDAO) DeleteByKey (m string) {
  dao.db.Where(&StringModel{ Key : m }).Delete(&StringModel{})
}

// EditByKey will edit all records in database
// with the same Key as parameter given
// using model given by parameter
func (dao *StringModelDAO) EditByKey (m string, newVals *StringModel) {
  dao.db.Table("string_models").Where(&StringModel{ Key : m }).Updates(newVals)
}

// SetKey will set Key
// to a value given by parameter
func (dao *StringModelDAO) SetKey (m *StringModel, newVal string) *StringModel {
  m.Key = newVal
  record := dao.ReadByID((m.ID))

  dao.db.Model(&record).Updates(m)

  return record
}

// ReadByValue will find all records
// matching the value given by parameter
func (dao *StringModelDAO) ReadByValue (m string) []StringModel {
  retVal := []StringModel{}
  dao.db.Where(&StringModel{ Value : m }).Find(&retVal)

  return retVal
}

// DeleteByValue deletes all records in database with
// Value the same as parameter given
func (dao *StringModelDAO) DeleteByValue (m string) {
  dao.db.Where(&StringModel{ Value : m }).Delete(&StringModel{})
}

// EditByValue will edit all records in database
// with the same Value as parameter given
// using model given by parameter
func (dao *StringModelDAO) EditByValue (m string, newVals *StringModel) {
  dao.db.Table("string_models").Where(&StringModel{ Value : m }).Updates(newVals)
}

// SetValue will set Value
// to a value given by parameter
func (dao *StringModelDAO) SetValue (m *StringModel, newVal string) *StringModel {
  m.Value = newVal
  record := dao.ReadByID((m.ID))

  dao.db.Model(&record).Updates(m)

  return record
}


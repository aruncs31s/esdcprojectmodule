package repository

import (
	commonModules "github.com/aruncs31s/esdcmodels"
)

// GetUserNotifications retrieves all notifications for a user
func (r *projectRepositoryReader) GetUserNotifications(userID uint, limit, offset int) ([]commonModules.Notification, error) {
	var notifications []commonModules.Notification
	err := r.db.
		Preload("TriggeredByUser").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetUnreadNotificationCount returns the count of unread notifications for a user
func (r *projectRepositoryReader) GetUnreadNotificationCount(userID uint) (int, error) {
	var count int64
	err := r.db.Model(&commonModules.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return int(count), err
}

// CreateNotification saves a new notification
func (w *projectRepositoryWriter) CreateNotification(notification *commonModules.Notification) error {
	return w.db.Create(notification).Error
}

// MarkNotificationAsRead marks a specific notification as read
func (w *projectRepositoryWriter) MarkNotificationAsRead(notificationID uint) error {
	return w.db.Model(&commonModules.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

// MarkAllNotificationsAsRead marks all notifications for a user as read
func (w *projectRepositoryWriter) MarkAllNotificationsAsRead(userID uint) error {
	return w.db.Model(&commonModules.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

// DeleteNotification removes a notification
func (w *projectRepositoryWriter) DeleteNotification(notificationID uint) error {
	return w.db.Delete(&commonModules.Notification{}, notificationID).Error
}

package service

//func (m *TgBot) GetOrderList(log string) []*model.UserOrderRes {
//	m.Mu.RLock()
//	user, ok := m.UserCH[log]
//	m.Mu.RUnlock()
//	if !ok {
//		m.Lg.Error("GetOrderList.err - пользователь: " + log + " не найдден")
//		return nil
//	}
//
//	orders := []*model.UserOrderRes{}
//
//	if len(user.OrderList) < 1 {
//		return orders
//	}
//
//	sort.Slice(orders, func(i, j int) bool {
//		return user.OrderList[i].Created.Before(user.OrderList[j].Created)
//	})
//
//	for _, order := range user.OrderList {
//		o := strconv.Itoa(order.OrderID)
//
//		accrual, err := decimal.NewFromString(order.Accrual)
//		if err != nil {
//			m.Lg.Error("GetOrderList.err - ошибка преобразования: " + err.Error())
//		}
//		accrualFloat64, _ := accrual.Float64()
//
//		orders = append(orders, &model.UserOrderRes{
//			Order:   o,
//			Status:  order.Status,
//			Accrual: float32(accrualFloat64),
//			Created: order.Created.Format(time.RFC3339),
//		})
//
//	}
//
//	return orders
//}
//
//func (m *TgBot) WithdrawnBalance(log, order string, amount decimal.Decimal) error {
//	m.Lg.Info("checkStatus.set.plus - установка новой транзакции для пользователя: " + log + " сумма: " + fmt.Sprintf("%v", amount) + " в статусе: processed")
//	err := m.Repo.SetTransaction(m.Ctx, log, order, "minus", amount)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (m *TgBot) GetInfoWithdrawnBalance(log string) ([]*model.WithdrawReq, error) {
//	list, err := m.Repo.GetInfoWithdrawnBalance(m.Ctx, log)
//	if err != nil {
//		return nil, err
//	}
//	return list, err
//}
//
//// SetOrder - формирование нового заказа
//func (m *TgBot) SetOrder(login string, order int) error {
//	m.Lg.Info(fmt.Sprintf("SetOrder.start - формирование нового заказа: %v", order))
//
//	create, err := m.Repo.CreateOrder(m.Ctx, order, login)
//	if err != nil {
//		m.Lg.Error(fmt.Sprintf("SetOrder.error - ошибка при добавлении заказа в БД %v", order))
//		return err
//	}
//
//	m.Mu.RLock()
//	user, ok := m.UserCH[login]
//	m.Mu.RUnlock()
//	if !ok {
//		m.Lg.Error("GetOrderList.err - пользователь: " + login + " не найдден")
//		return nil
//	}
//	user.OrderList[order] = &model.Order{
//		OrderID: order,
//		Status:  model.NEW,
//		Created: create,
//	}
//	m.Lg.Info(fmt.Sprintf("SetOrder.progress - заказ: %v, успешно добавлен в БД и Кеш", order))
//
//	go m.processGetStatus(m.Ctx, login, order, m.paramGetStatus)
//	return nil
//}
//
//// checkStatus - мапинг статусов
//func (m *TgBot) checkStatus(ctx context.Context, res *model.AccrualRes, order int, user string) (bool, error) {
//	switch res.Status {
//	case model.PROCESSED:
//		err := m.Repo.SetBonus(ctx, order, model.PROCESSED, fmt.Sprintf("%v", res.Accrual))
//		if err != nil {
//			return false, err
//		}
//
//		dec, err := decimal.NewFromString(res.Accrual.String())
//		if err != nil {
//			return false, fmt.Errorf("неверный формат accrual: %v", res.Accrual)
//		}
//		err = m.Repo.SetTransaction(ctx, user, fmt.Sprintf("%v", order), "plus", dec)
//		if err != nil {
//			return false, err
//		}
//		m.Lg.Info("checkStatus.set.plus - установка новой транзакции для пользователя: " + user + " сумма: " + fmt.Sprintf("%v", dec) + " в статусе: processed")
//
//		return true, nil
//	case model.PROCESSING:
//		err := m.Repo.SetStatus(ctx, order, model.PROCESSING)
//		if err != nil {
//			return false, err
//		}
//		return false, nil
//	case model.INVALID:
//		err := m.Repo.SetStatus(ctx, order, model.INVALID)
//		if err != nil {
//			return false, err
//		}
//		return true, nil
//	}
//	return false, nil
//}
//
